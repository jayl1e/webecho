# why
- test your reverse proxy deployment

# how
- just kubectl apply -f deploy.yaml
- use curl or browser

# functions
- it just echo your request
- use param pre_sleep=90s&post_sleep=90s to test your reverse proxy timeout before receiving header and after receiving header
- use skip_body=1 param or Skip-Body: 1 header to stop webecho return your body
- upload a big body to test reverse proxy, webecho will just stream it back
