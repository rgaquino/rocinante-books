import csv
import json
import re
from datetime import datetime


def format_title(title: str) -> (str, str):
    colon = row[0].find(":")
    if colon:
        sub = title.title()[colon+1:].strip()
    title = title.title()[:colon].strip()
    return title, sub


def format_author(author: str) -> str:
    authors = author.split(",")
    author = authors[0]
    if len(authors) == 1:
        return author
    count = 1
    for a in authors[1:]:
        a = a.strip()
        if count == len(authors)-1:
            if len(authors) > 2:
                author += f", and {a}"
            else:
                author += f" and {a}"
        else:
            author += f", {a}"
        count += 1

    return author


def format_highlight(highlight: str) -> str:
    highlight = highlight.strip()
    if highlight[0].islower():
        highlight = f"[{highlight[0].upper()}]{highlight[1:]}"
    if highlight[-1].isalnum():
        highlight = f"{highlight}."
    return highlight


def format_slug(title: str, author: str) -> str:
    title = re.sub(r'[^A-Za-z0-9 ]+', '', title)
    author = re.sub(r'[^A-Za-z0-9 ]+', '', author)
    return '-'.join(
        word for word in f'{title} {author}'.lower().split(' '))


BOOK_TITLE = "title"
BOOK_SUBTITLE = "subtitle"
BOOK_AUTHOR = "author"
BOOK_HIGHLIGHTS = "highlights"
BOOK_SLUG = "slug"
BOOK_FINISHED_AT = "finishedAt"
BOOK_LAST_FINISHED_AT = "lastFinishedAt"
BOOK_CATEGORY = "category"

now: str = datetime.now().isoformat()
book: dict = {
    BOOK_TITLE: "",
    BOOK_AUTHOR: "",
    BOOK_HIGHLIGHTS: [],
    BOOK_SLUG: "",
    BOOK_FINISHED_AT: [now],
    BOOK_LAST_FINISHED_AT: now,
    BOOK_CATEGORY: ""
}

with open('sample_in.csv', 'r') as file:
    reader = csv.reader(file, delimiter=',')
    count = 0
    for row in reader:
        if count == 1:
            book[BOOK_TITLE], book[BOOK_SUBTITLE] = format_title(row[0])
        elif count == 2:
            book[BOOK_AUTHOR] = format_author(row[0][3:])
        elif count < 8:
            count += 1
            continue
        else:
            highlight = format_highlight(row[3])
            book[BOOK_HIGHLIGHTS].append(highlight)
        count += 1

    book[BOOK_SLUG] = format_slug(book[BOOK_TITLE], book[BOOK_AUTHOR])
    print(json.dumps(book, indent=4))
