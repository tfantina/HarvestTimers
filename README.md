# ⏲️ Harvest Timers

Because I contract with multiple clients I have two Harvest accounts, one for the agency
and one for myself. This makes it cumbersome to know how many hours I've worked for the
month against my goal of ~140. 

This takes credentials for both accounts and queries the `reports/time/team` endpoint 
to get my total hours and sums them.

It's rudimentary but it does exactly what I need. While I could make it more flexible
(custom setup options, additional Harvest accounts, additional data etc.), this was 
built just for me and I have too much on the go to develop features I don't need.

## 🏗️ Setup

Create a `.env` file in the same directory as the project.
Include 4 variables:
`ACCESS_TOKEN`, `ACCOUNT_ID` and `ACCESS_TOKEN2`, `ACCOUNT_ID2` which correspond to 
two separate Harvest accounts.


`go run main.go`

## 🎶 Additional Notes & AI Disclaimer

This is my first project in Go. I loathe AI but I've been experimenting with OpenCode. 
With my schedule and commitments; I would not have been able to get this off the ground 
in Go without AI.

OpenCode generated the project structure and most of the functions in `main.go` the 
prompts were based on a single account because I just wanted to get a good Go pattern 
for making API calls and unmarshaling JSON.  From there I modified/created requisite 
types, specifically:
 * `TimeReport`
 * `TeamMember`

I added additional functionality to query a second Harvest account and store environment
variables in a .env file.
