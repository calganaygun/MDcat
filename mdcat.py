from sys import argv
from requests import request
from os.path import split

MD_FILE_PATH = argv[1]
MD_FILE_DIR, MD_FILE_NAME = split(MD_FILE_PATH)

with open(MD_FILE_PATH, "r", encoding="utf-8") as md_file:
    MD_FILE_CONTENT = md_file.read()

url = "https://api.github.com/markdown"

headers = {
  'Accept': 'application/vnd.github.v3+json',
  'Content-Type': 'application/json'
}

payload = {"text": MD_FILE_CONTENT}
response = request("POST", url, headers=headers, json=payload)

if response.status_code == 200:
    with open("template.html", "r", encoding="utf-8") as template_file:
        TEMPLATE_CONTENT = template_file.read()
    html_content = TEMPLATE_CONTENT.replace("$MD_TITLE", MD_FILE_NAME).replace("$MD_HTML", response.text)
    with open(MD_FILE_DIR + MD_FILE_NAME.strip(".md") + ".html", "w", encoding="utf-8") as export_file:
        export_file.write(html_content)
else:
    print("ERROR:", response.json()["message"])
