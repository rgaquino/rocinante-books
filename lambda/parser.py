import csv
import json
import re
import sys
import requests
from bs4 import BeautifulSoup
from datetime import datetime
import urllib.request

INVALID_FINAL_CHARS = [",", ":", ";"]


def format_title_author(title_author: str) -> (str, str, str):
    author_start = title_author.rfind("(")
    author_end = title_author.rfind(")")
    author = format_author(title_author[author_start + 1: author_end])
    title, sub = format_title(title_author[:author_start])
    return title, sub, author


def format_title(title: str) -> (str, str):
    colon = title.find(":")
    sub = ""
    if colon >= 0:
        sub = title.title()[colon + 1:].strip()
    else:
        colon = len(title)
    title = title.title()[:colon].strip()
    return title, sub


def format_author(author: str) -> str:
    authors = author.strip().split(",")
    author = authors[0]
    if len(authors) == 1:
        return author
    count = 1
    for a in authors[1:]:
        a = a.strip()
        if count == len(authors) - 1:
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
    elif highlight[-1] in INVALID_FINAL_CHARS:
        highlight = f"{highlight[:-1]}."
    return highlight


def format_category(category: str) -> str:
    return category.strip()


def format_slug(title: str, author: str) -> str:
    title = re.sub(r"[^A-Za-z0-9 ]+", "", title)
    author = re.sub(r"[^A-Za-z0-9 ]+", "", author)
    return "-".join(word for word in f"{title} {author}".lower().split(" "))


BOOK_ID = "id"
BOOK_TITLE = "title"
BOOK_SUBTITLE = "subtitle"
BOOK_AUTHOR = "author"
BOOK_HIGHLIGHTS = "highlights"
BOOK_SLUG = "slug"
BOOK_FINISHED_AT = "finishedAt"
BOOK_LAST_FINISHED_AT = "lastFinishedAt"
BOOK_CATEGORY = "category"


def new_book() -> dict:
    now: str = datetime.now().isoformat()
    return {
        BOOK_ID: 0,
        BOOK_TITLE: "",
        BOOK_AUTHOR: "",
        BOOK_HIGHLIGHTS: [],
        BOOK_SLUG: "",
        BOOK_FINISHED_AT: [now],
        BOOK_LAST_FINISHED_AT: now,
        BOOK_CATEGORY: "",
    }


def parse_kindle(fn) -> dict:
    book = new_book()
    with open(fn, "r") as file:
        reader = csv.reader(file, delimiter=",")
        count = 0
        for row in reader:
            if count == 1:
                book[BOOK_TITLE], sub = format_title(row[0])
                if sub:
                    book[BOOK_SUBTITLE] = sub
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
    return book


def parse_manual(fn) -> dict:
    book = new_book()
    file = open(fn, "r")
    count = 0
    for line in file:
        if count == 0:
            book[BOOK_TITLE], sub = format_title(line)
            if sub:
                book[BOOK_SUBTITLE] = sub
        elif count == 1:
            book[BOOK_AUTHOR] = format_author(line)
        elif count == 2:
            book[BOOK_CATEGORY] = format_category(line)
        else:
            if line.strip():
                highlight = format_highlight(line)
                book[BOOK_HIGHLIGHTS].append(highlight)
        count += 1
    book[BOOK_SLUG] = format_slug(book[BOOK_TITLE], book[BOOK_AUTHOR])
    return book


def parse_clip(fn) -> dict:
    book = new_book()
    file = open(fn, "r")
    lines = []
    for line in file:
        lines.append(line)

    book[BOOK_TITLE], sub, book[BOOK_AUTHOR] = format_title_author(lines[0])
    if sub:
        book[BOOK_SUBTITLE] = sub

    count = 4
    while count < len(lines):
        highlight = format_highlight(lines[count - 1])
        book[BOOK_HIGHLIGHTS].append(highlight)
        count += 5
    book[BOOK_SLUG] = format_slug(book[BOOK_TITLE], book[BOOK_AUTHOR])
    return book


def get_image_url(title) -> str:
    title = title.replace(" ", "+")
    response = requests.get(
        f"https://www.goodreads.com/search/index.xml?key=NMudcgJgw9gQPDRDb4jeiQ&q={title}&search[field]=title"
    )
    if response.status_code == 200:
        decoded = BeautifulSoup(response.text, features="html.parser")
        image_url = decoded.goodreadsresponse.search.results.findAll(
            "work")[0].best_book.image_url.getText()
        image_url = image_url[0:image_url.rfind(
            "._")] + image_url[image_url.rfind("_.")+1:]
        return image_url
    else:
        return ''


def main():
    book = {}

    if len(sys.argv) < 3:
        print("Usage: p3 parser.py <source_type> <file_name>")
    elif sys.argv[1] == "kindle":
        book = parse_kindle(sys.argv[2])
    elif sys.argv[1] == "manual":
        book = parse_manual(sys.argv[2])
    elif sys.argv[1] == "clip":
        book = parse_clip(sys.argv[2])
    else:
        print("<source_type> can only be either 'kindle' or 'manual'")

    if len(book) > 0:
        # print(json.dumps(book, indent=4))
        image_url = get_image_url(book[BOOK_TITLE])
        urllib.request.urlretrieve(image_url, "sample.jpg")


if __name__ == "__main__":
    main()
