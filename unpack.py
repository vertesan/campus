from pathlib import Path
from rich.console import Console
import UnityPy
import UnityPy.config
import json

ASSETBUNDLE_DIR = "cache/assets"
IMG_DIR = "cache/img"
ESRGAN_DIR = "cache/esrgan"
DOWNLOADED_FILE_PATH = "cache/octo_downloaded.json"
UnityPy.config.FALLBACK_UNITY_VERSION = "2022.3.21f1"
console = Console()


def info(msg: str):
    console.print(f"[bold blue]>>> [Info][/bold blue] {msg}")


def succeed(msg: str):
    console.print(f"[bold green]>>> [Succeed][/bold green] {msg}")


def error(msg: str):
    console.print(f"[bold red]>>> [Error][/bold red] {msg}")


def warn(msg: str):
    console.print(f"[bold yellow]>>> [Warning][/bold yellow] {msg}")


def unpack_to_image(asset_bytes: bytes, dest: str):
    env = UnityPy.load(asset_bytes)
    for obj in env.objects:
        if obj.type.name == "Texture2D":
            try:
                data = obj.read()
                # one of the QA employees messed up upper and lower case of assetname,
                # traditionally they are all written in lowercase
                if data.name == "img_general_icon_exam-effect_examItemfirelimitadd":
                    filename = data.name.lower()
                else:
                    filename = data.name
                dest_path = Path(dest, filename).with_suffix(".png")
                dest_path.parent.mkdir(exist_ok=True)
                img = data.image
                img.save(dest_path)
                info(f"Converted '{filename}' to png.")
            except:
                error(f"Failed to convert '{filename}' to image.")


def unpack_action(octo_diff: dict[str, str]):
    for name, _ in octo_diff.items():
        try:
            raw = Path(ASSETBUNDLE_DIR, name).read_bytes()
            if raw[:5] != b"Unity":
                warn(f"'{name}' is not a unity asset, skip processing.")
                continue
            unpack_to_image(raw, IMG_DIR)
        except:
            error(f"Failed to process '{name}'.")


def scale_with_esrgan(octo_diff: dict[str, str]):
    from esrgan import convert_one
    for name, _ in octo_diff.items():
        if name.startswith("img_general_cidol-") and name.endswith("-full"):
            convert_one(
                str(Path(IMG_DIR, name + ".png")),
                ESRGAN_DIR,
                extension="webp",
                to_size=True,
                c_size=(1440, 2560),
            )
        if name.startswith("img_general_csprt-") and name.endswith("_full"):
            convert_one(
                str(Path(IMG_DIR, name + ".png")),
                ESRGAN_DIR,
                extension="webp",
                to_size=True,
                c_size=(2560, 1440),
            )
        if name.startswith("img_adv_still_"):
            convert_one(
                str(Path(IMG_DIR, name + ".png")),
                ESRGAN_DIR,
                extension="webp",
                to_size=True,
                c_size=(1440, 2560),
            )


def main():
    Path(IMG_DIR).mkdir(exist_ok=True)
    Path(ESRGAN_DIR).mkdir(exist_ok=True)
    with open(DOWNLOADED_FILE_PATH) as fp:
        octo_diff: dict[str, str] = json.load(fp)
    unpack_action(octo_diff)
    scale_with_esrgan(octo_diff)


if __name__ == "__main__":
    main()
