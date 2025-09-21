# RSSHub

# Description
This is a service that collects publications from various sources that provide RSS feeds (news sites, blogs, forums). It helps users stay informed in one place without the need to visit each website manually.
Such a tool is useful for journalists, researchers, analysts, and anyone who wants to stay updated on topics of interest without unnecessary noise. This kind of application makes information more accessible and centralized.

# Installation
docker compose up --build
make

# Usage
./rsshub --help

./rsshub fetch

./rsshub add --name "tech-crunch" --url "https://techcrunch.com/feed/"

./rsshub set-interval --duration 5m

./rsshub set-workers --count 5

./rsshub list --num 5

./rsshub delete --name "tech-crunch"

./rsshub articles --feed-name "tech-crunch" --num 5