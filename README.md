# Introduction

This is aggreGator-go - a multi user tool for aggregating and viewing rss feeds from the command line written in Go. This project was written as part of an exercise for the learning program boot.dev.

## Requirements

1. The latest [Go toolchain](https://golang.org/dl/)
2. A local PostGres database

## Installation

Installation can be done in the standard manner - clone the repository, then run:

```bash
go install
```

In the root of the repo.

## Config

Create a `.gatorconfig.json` file in your home directory containing the following json structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your postgres database string

## Usage

Create a user:

```bash
aggreGator-go register <name>
```

Add an RSS feed:

```bash
aggreGator-go addfeed <name> <url>
```

Start the aggregator (please be responsible when setting the timestep - you don't want to accidentally DOS the feeds you are subscribed to):

```bash
aggreGator-go agg <timestep>
```

Veiw the posts:

```bash
aggreGator-go browse <limit>
```

Other commands are available, but these are left as an exercise to the reader (JK, there's a help command coming soon (tm))
