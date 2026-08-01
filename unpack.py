import os
from pathlib import Path
from rich.console import Console
import UnityPy
import UnityPy.config
import json

ASSETBUNDLE_DIR = "cache/assets"
IMG_DIR = "cache/img"
ESRGAN_DIR = "cache/esrgan"
DOWNLOADED_FILE_PATH = "cache/octo_downloaded.json"
UnityPy.config.FALLBACK_UNITY_VERSION = "6000.0.77f1"
console = Console()


def info(msg: str):
    console.print(f"[bold blue]>>> [Info][/bold blue] {msg}")


def succeed(msg: str):
    console.print(f"[bold green]>>> [Succeed][/bold green] {msg}")


def error(msg: str):
    console.print(f"[bold red]>>> [Error][/bold red] {msg}")


def warn(msg: str):
    console.print(f"[bold yellow]>>> [Warning][/bold yellow] {msg}")


def unpack_to_image(objPath: str, dest_dir: str):
    env = UnityPy.load(objPath)
    for obj in env.objects:
        if obj.type.name in ["Texture2D", "Sprite"]:
            try:
                data = obj.parse_as_object()
                dest_path = os.path.join(dest_dir, data.m_Name)
                dest_path, ext = os.path.splitext(dest_path)
                dest_path = dest_path + ".png"
                img = data.image
                img.save(dest_path)
                info(f"Converted '{data.m_Name}' to png.")
            except Exception as e:
                error(f"Failed to convert '{data.m_Name}' to image.")
                error(f"Msg: {e}")


def unpack_action(octo_diff: dict[str, str]):
    for name, _ in octo_diff.items():
        if name.endswith(".txt") or name.endswith(".acb") or name.endswith(".acf"):
            continue
        try:
            objPath = os.path.join(ASSETBUNDLE_DIR, name)
            with open(objPath, "rb") as file:
                sig = file.read(5)
                if sig != b"Unity":
                    warn(f"'{name}' is not a unity asset, skip processing.")
                    continue
            unpack_to_image(objPath, IMG_DIR)
        except Exception as e:
            error(f"Failed to process '{name}'.")
            error(f"Msg: {e}")


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
