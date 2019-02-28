# go-resume

go-resume is a command line tool for your JSON Resume

## What is JSON Resume

*json resume* is JSON-Schema that is used here to define and validate our proposed resume

https://github.com/jsonresume/resume-schema

## Usage

### validation your JSON-Resume

```bash
resume validate [your resume file path]
```

i.g.
```bash
resume validate $GOPATH/github.com/usk81/tests/resume.json
```

### Export (HTML)

```bash
resume export [your template directory path] [output directory path] [your resume file path]
```

i.g.
```bash
resume validate $GOPATH/github.com/usk81/themes/html ~/ $GOPATH/github.com/usk81/tests/resume.json 
```

## Installation

```bash
go get -u github.com/usk81/resume
```

## Features

- [x] validates a JSON-Resume file
- export
  - [x] use custom HTML template
     - golang HTML template (html/template) ONLY
  - [ ] exports PDF
  - [ ] exports MS Excel (XLSX)