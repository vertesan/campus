.\campus.exe --ab --db
if (Test-Path "cache/newab_flag") {
  python unpack.py
}

pause
