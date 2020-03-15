import csv
import json
import re
from datetime import datetime

now = datetime.now().isoformat()

book = {
    "title": "",
    "author": "",
    "highlights": [],
    "finishedAt": [now],
    "lastFinishedAt": now,
    "category": "",
    "slug": ""
}

with open('sample_in.csv', 'r') as file:
    reader = csv.reader(file, delimiter=',')
    count = 0
    for row in reader:
        if count == 1:
            colon = row[0].find(":")
            if colon:
                book["subtitle"] = row[0].title()[colon+1:].strip()
            book["title"] = row[0].title()[:colon].strip()
        elif count == 2:
            book["author"] = row[0].title()[3:]
        elif count < 8:
            count += 1
            continue
        else:
            book["highlights"].append(row[3])
        count += 1

    title = re.sub(r'\W+ ', '', book["title"])
    author = re.sub(r'\W+ ', '', book["author"])
    book["slug"] = '-'.join(
        word for word in f'{title} {author}'.lower().split(' '))

    print(json.dumps(book, indent=4))
