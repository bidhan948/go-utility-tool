
# 🧹 cln — Tiny Disk Cleanup CLI

Clean your disk, find duplicates, and remove junk — all from your terminal ✨  

---

## ❓ What is `cln`?

A lightweight Go-based tool to:
- 🕵️ Scan large files  
- 🔁 Find duplicate files  
- 💡 Suggest junk folders  
- 🗑️ Remove unwanted files  

---

## ⚙️ Installation

```bash
git clone https://github.com/bidhan948/go-utility-tool.git
cd go-utility-tool/cln
go build -o cln ./cmd/cln
````

Run it anywhere:

```bash
./cln --help
```

---

## 💬 Common Questions

### 💡 Q: How do I **scan** for big files?

```bash
./cln scan /path/to/folder --min 100MB
```

🧾 Lists files **≥ 100 MB**.
If you skip `--min`, it shows **all files**.

Example:

```bash
./cln scan ~/Downloads
./cln scan ~/Videos --min 1GB
```

---

### 🔁 Q: How do I **find duplicates**?

```bash
./cln dup /path/to/folder --min 10MB
```

Finds duplicate files with same hash and size.
🧱 Use `--json` for structured output.

---

### 💡 Q: How do I get **junk suggestions**?

```bash
./cln suggest
```

Shows common cache/log/temp folders that may be safely removed.
Use `--json` to get JSON output for scripting.

---

### 🗑️ Q: How do I **delete files**?

```bash
./cln rm file1 file2 ...
```

🧠 Safety first:

* Prompts for confirmation.
* Add `--force` to skip confirmation.
* Add `--dry` for a dry run (no deletion).

Example:

```bash
./cln rm ~/Downloads/bigfile.iso
./cln rm ~/Downloads/*.tmp --force
```

---

### 🧾 Q: Want **JSON** output?

```bash
./cln scan ~/Downloads --json
```

Perfect for automation, scripting, or piping to tools like `jq`.

---

### 🧭 Q: Need to check version?

```bash
./cln version
```

---

## 🧰 Quick Reference

| Command      | Description             | Example                            |
| ------------ | ----------------------- | ---------------------------------- |
| 🕵️ `scan`   | Scan for large files    | `cln scan ~/Downloads --min 500MB` |
| 🔁 `dup`     | Find duplicate files    | `cln dup ~/Videos --min 50MB`      |
| 💡 `suggest` | Show junk suggestions   | `cln suggest`                      |
| 🗑️ `rm`     | Delete files safely     | `cln rm old.log --force`           |
| 🧾 `--json`  | JSON output for scripts | `cln scan --json`                  |
| ⚙️ `version` | Show version            | `cln version`                      |

---

## 🌈 Tips

* Combine with `grep` or `jq` for automation

  ```bash
  cln scan ~/Downloads --json | jq '.[] | select(.size > 100000000)'
  ```
* Use `--min` smartly for quick filtering.
* Works beautifully with colored terminals.

---

## 🪄 Example Run

```bash
./cln scan ~/Downloads --min 50MB
```

Output:

```
📦 SIZE      🕓 MODIFIED             📄 PATH
-----------------------------------------------
123 MiB      2025-11-09 14:05:21    🗎 /home/user/Downloads/movie.mp4
98 MiB       2025-11-08 11:45:10    🗎 /home/user/Downloads/archive.zip

✅ Listed 2 file(s) ≥ 50MB — total 221 MiB
```

---

### 🧡 Enjoy a cleaner disk with `cln`!