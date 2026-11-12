package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"vertesan/campus/crypto"
	"vertesan/campus/master"
	"vertesan/campus/network/rpc"
	"vertesan/campus/octo"
	"vertesan/campus/proto/mapping"
	octopb "vertesan/campus/proto/octo"
	"vertesan/campus/proto/pcommon"
	"vertesan/campus/proto/papi"
	"vertesan/campus/proto/papicommon"
	"vertesan/campus/proto/pmaster"
	"vertesan/campus/proto/penum"

	"github.com/goccy/go-yaml"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestEmptyPayloadCryptoAndRpcRoundtrip(t *testing.T) {
	key := bytes.Repeat([]byte{0x42}, 16)
	iv := bytes.Repeat([]byte{0x24}, 16)

	var enc bytes.Buffer
	crypto.Encrypt(key, iv, bytes.NewReader(nil), &enc)
	if enc.Len() != aes.BlockSize {
		t.Fatalf("expected 16-byte padded block for empty plaintext, got %d", enc.Len())
	}

	var dec bytes.Buffer
	crypto.Decrypt(key, iv, bytes.NewReader(enc.Bytes()), &dec)
	if dec.Len() != 0 {
		t.Fatalf("expected empty plaintext after Decrypt, got %d bytes", dec.Len())
	}

	// Default protobuf message serializes to 0 bytes; QuaCodec Serialize/Deserialize must roundtrip cleanly.
	emptyRespBytes, err := proto.Marshal(&papi.Empty{})
	if err != nil {
		t.Fatalf("marshal Empty: %v", err)
	}
	wire := rpc.Serialize(emptyRespBytes)
	got := rpc.Deserialize(wire)
	if len(got) != 0 {
		t.Fatalf("expected empty deserialized bytes, got %d", len(got))
	}
}

func TestDecryptOctoListWithoutContentLength(t *testing.T) {
	orig := &octopb.Database{
		Revision:  86,
		UrlFormat: "https://example.invalid/{o}",
		AssetBundleList: []*octopb.Data{
			{Name: "img_general_test", ObjectName: "obj1", Md5: "0123456789abcdef0123456789abcdef"},
		},
	}
	plain, err := proto.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal octo db: %v", err)
	}

	key := sha256.Sum256([]byte(octo.OCTO_API_KEY))
	iv := bytes.Repeat([]byte{0xAB}, aes.BlockSize)
	var encBuf bytes.Buffer
	crypto.Encrypt(key[:], iv, bytes.NewReader(plain), &encBuf)

	// Simulate unknown Content-Length (-1) and a 1-byte prefix offset + 3-byte gap.
	stream := append([]byte{0x01}, iv...)
	stream = append(stream, 0xFE, 0xFD, 0xFC)
	stream = append(stream, encBuf.Bytes()...)

	got, err := octo.DecryptOctoList(bytes.NewReader(stream), 1, -1)
	if err != nil {
		t.Fatalf("DecryptOctoList failed with contentLen=-1: %v", err)
	}
	if got.Revision != 86 || len(got.AssetBundleList) != 1 || got.AssetBundleList[0].Name != "img_general_test" {
		t.Fatalf("unexpected decrypted octo db: %+v", got)
	}
	_ = cipher.BlockMode(nil)
}

func TestDeobfuscateSmallAndStandardAssetBundle(t *testing.T) {
	abName := "img_general_skillcard_test"

	// 1. Short (<256 bytes) payload must not panic and must roundtrip cleanly.
	shortOrig := []byte("UnityFS\x00short-bundle-payload")
	obfShort := octo.Deobfuscate(bytes.NewReader(bytes.Clone(shortOrig)), abName)
	if bytes.Equal(obfShort[:5], []byte("Unity")) {
		t.Fatalf("expected obfuscated header to differ from Unity")
	}
	restoredShort := octo.Deobfuscate(bytes.NewReader(obfShort), abName)
	if !bytes.Equal(restoredShort, shortOrig) {
		t.Fatalf("short bundle roundtrip mismatch: got %q, want %q", restoredShort, shortOrig)
	}

	// 2. Standard (>256 bytes) payload only masks first 256 bytes.
	longOrig := append([]byte("UnityFS\x00"), bytes.Repeat([]byte{0x5A}, 400)...)
	obfLong := octo.Deobfuscate(bytes.NewReader(bytes.Clone(longOrig)), abName)
	if !bytes.Equal(obfLong[256:], longOrig[256:]) {
		t.Fatalf("bytes past 256 should remain untouched")
	}
	restoredLong := octo.Deobfuscate(bytes.NewReader(obfLong), abName)
	if !bytes.Equal(restoredLong, longOrig) {
		t.Fatalf("long bundle roundtrip mismatch")
	}
}
func TestProto360SchemasAndSQLCipherMasterDecrypt(t *testing.T) {
	required360Tables := []string{
		"Master.RefineAdv",
		"Master.RefineDefaultMemory",
		"Master.RefinePark",
		"Master.RefineProduceCardCustomizeGroup",
		"Master.RefineSetting",
		"Master.RefineStartMotion",
		"Master.RefineStepAudition",
		"Master.RefineStepAuditionMotion",
		"Master.RefineStepSetupCost",
		"Master.ExamMainBuff",
		"Master.CharacterPosition",
		"Master.ProduceCardSkin",
		"Master.CostumePhotoGroup",
		"Master.ProduceNextIdolAuditionMasterRankingSeason",
		"Master.ResearchMemoryRerollCost",
		"Master.ExamSetting",
		"Master.Setting",
		"Common.ProduceDescriptionSegment",
	}
	for _, name := range required360Tables {
		if _, ok := mapping.ProtoMap[name]; !ok {
			t.Fatalf("missing 3.6.0 entry %q in mapping.ProtoMap", name)
		}
	}

	// Check 3.6.0 Plan4 / Refine fields & enums compile and populate cleanly.
	homeResp := &papi.HomeEnterResponse{
		Refines: []*papicommon.Refine{
			{RefineId: "refine-001"},
		},
	}
	if len(homeResp.GetRefines()) != 1 || homeResp.GetRefines()[0].GetRefineId() != "refine-001" {
		t.Fatalf("unexpected HomeEnterResponse.Refines: %+v", homeResp)
	}

	// 2. Test full SQLCipher encrypted SQLite -> master.DecryptAll roundtrip in an isolated temp cwd.
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.MkdirAll(master.MASTER_RAW_PATH, 0755); err != nil {
		t.Fatal(err)
	}

	hexKey := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	dbPath := filepath.Join(master.MASTER_RAW_PATH, "Setting")
	db, err := sql.Open("sqlite3", dbPath+"?_pragma_key=x'"+hexKey+"'")
	if err != nil {
		t.Fatalf("open sqlcipher db: %v", err)
	}
	if _, err := db.Exec("CREATE TABLE Setting (data BLOB);"); err != nil {
		db.Close()
		t.Fatalf("create table: %v", err)
	}

	settingRow := &pmaster.Setting{
		Id:                      "setting-1",
		Plan4ViewConditionSetID: "cond-plan4-unlock",
	}
	rowBytes, err := proto.Marshal(settingRow)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO Setting (data) VALUES (?);", rowBytes); err != nil {
		db.Close()
		t.Fatal(err)
	}
	db.Close()

	// Also verify ProduceCard with 3.6.0 Plan4 enum & ProduceDescriptionSegment originMove fields
	seg := &pcommon.ProduceDescriptionSegment{
		ProduceCardMovePositionType:    penum.ProduceCardMovePositionType_ProduceCardMovePositionType_Remove,
		OriginMoveProduceExamTriggerId: "trig-plan4-001",
		OriginMoveProduceExamEffectId:  "eff-plan4-001",
	}
	if seg.GetProduceCardMovePositionType() != penum.ProduceCardMovePositionType_ProduceCardMovePositionType_Remove ||
		seg.GetOriginMoveProduceExamEffectId() != "eff-plan4-001" {
		t.Fatalf("unexpected 3.6.0 ProduceDescriptionSegment fields: %+v", seg)
	}
	resp := &papi.MasterGetResponse{
		MasterTag: &papi.MasterGetResponse_MasterTag{
			Version: "3.6.0-test",
			MasterTagPacks: []*papi.MasterGetResponse_MasterTagPack{
				{Type: "Setting", CryptoKey: hexKey},
			},
		},
	}
	master.DecryptAll(resp, false)

	yamlOut, err := os.ReadFile(filepath.Join(master.MASTER_YAML_PATH, "Setting.yaml"))
	if err != nil {
		t.Fatalf("read decrypted Setting.yaml: %v", err)
	}
	jsonFromYaml, err := yaml.YAMLToJSON(yamlOut)
	if err != nil {
		t.Fatalf("YAMLToJSON on Setting.yaml: %v", err)
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(jsonFromYaml, &arr); err != nil || len(arr) != 1 {
		t.Fatalf("unexpected Setting.yaml array: %v, len=%d", err, len(arr))
	}
	var gotSetting pmaster.Setting
	if err := (protojson.UnmarshalOptions{DiscardUnknown: false}).Unmarshal(arr[0], &gotSetting); err != nil {
		t.Fatalf("unmarshal decrypted Setting: %v", err)
	}
	if gotSetting.Plan4ViewConditionSetID != "cond-plan4-unlock" {
		t.Fatalf("got Plan4ViewConditionSetID=%q, want cond-plan4-unlock", gotSetting.Plan4ViewConditionSetID)
	}
}
