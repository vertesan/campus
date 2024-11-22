import re
import os
from PIL import Image
import cloudinary
import cloudinary.uploader
from pathlib import Path
import json
from unpack import unpack_to_image

ASSET_DIR = "cache/assets"
UPLOAD_IMG_DIR = "cache/webimg"
OCTO_RECORD_PATH = "cache/octo_record.json"
UPLOAD_RECORD_PATH = "cache/uploaded_record.json"
upload_reg_list: dict[str, tuple] = {
    r"^img_general_csprt-\d-\d{4}_full$": (1820, 1024),
    r"^img_general_cidol-\w+-\d-\d{3}_\d-full$": (1024, 1820),
    r"^img_general_icon_exam-\w+$": (128, 128),
    r"^img_general_icon_produce-.+$": (128, 128), # include effect icons
    r"^img_general_pitem_\d-\d+$": (256, 256),
    r"^img_general_skillcard_.+": (256, 256),
    r"^img_sd_[a-z]{4}_face-00": (256, 256),
    r"^img_general_event_.+banner$": (745, 256),
    r"^img_general_achievement_.+$": (512, 512),
    r"^img_chr_.+full$": (2048, 2048),
    r"^img_chr_.+thumb-circle$": (256, 256),
    r"^img_general_sign_\w+_\d{2}$": (1024, 1024),
    r"^img_general_meishi_illust_stamp-.+$": (128, 128),
}

# Configuration
cloudinary.config(
    cloud_name=os.environ["CLOUDINARY_CLOUD_NAME"],
    api_key=os.environ["CLOUDINARY_API_KEY"],
    api_secret=os.environ["CLOUDINARY_API_SECRET"],
    secure=True,
)


def upload(name: str, pth: str):
    # Upload an image
    upload_result = cloudinary.uploader.upload(
        pth,
        public_id=name,
        folder="gkms",
        overwrite=True,
    )
    print(upload_result["secure_url"])


def convert_to_webp(src: str, dst: str, c_size: tuple, delete_src: bool):
    img = Image.open(src)
    img = img.resize(c_size, resample=Image.LANCZOS, reducing_gap=3)
    img.save(fp=dst, format="webp", lossless=False, quality=80)
    if delete_src and src != dst:
        img.close()
        Path(src).unlink(missing_ok=True)


def filter_assets(name: str) -> bool:
    for reg in upload_reg_list.keys():
        if re.match(reg, name) != None:
            return True
    return False


def get_csize(name: str) -> tuple:
    for reg, csize in upload_reg_list.items():
        if re.match(reg, name):
            return csize
    raise "No predefined csize to match the image, this should not happen."


def main():
    print("Start to run unpack_upload.py")
    upload_rcd: dict[str, str] = {}
    if Path(UPLOAD_RECORD_PATH).exists():
        with open(UPLOAD_RECORD_PATH) as fp:
            upload_rcd = json.load(fp)
    with open(OCTO_RECORD_PATH) as fp:
        octo_rcd: dict[str, str] = json.load(fp)

    diff_rcd = octo_rcd
    # filter MD5
    diff_rcd = dict(filter(lambda it: upload_rcd.get(it[0]) != it[1], diff_rcd.items()))
    # filter regex
    diff_rcd = dict(filter(lambda it: filter_assets(it[0]), diff_rcd.items()))

    for name, md5 in diff_rcd.items():
        raw = Path(ASSET_DIR, name).read_bytes()
        unpack_to_image(raw, UPLOAD_IMG_DIR)
        png_pth = Path(UPLOAD_IMG_DIR, name + ".png")
        webp_pth = Path(UPLOAD_IMG_DIR, name + ".webp")
        convert_to_webp(str(png_pth), str(webp_pth), get_csize(name), True)
        upload(name, str(webp_pth))
        upload_rcd[name] = md5
        Path(UPLOAD_RECORD_PATH).write_text(json.dumps(upload_rcd, indent=2))

    print("Process of unpack_upload.py completed.")


if __name__ == "__main__":
    main()
