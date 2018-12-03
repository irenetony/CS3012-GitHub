# CS3012-GitHub
This program will take a github username and display the all the users that username follows and the amount of repos they have.

## How to run 
- Insert your github access token into the "token" area in the .go file. 
- Run the .go file on termainal using `go run GitHubAccess.go`. 
- Open up `localhost:9000/web` on your browser. This should display the website and ask for a user name.

![Starting page](https://github.com/irenetony/CS3012-GitHub/raw/master/Screenshots/Start.png) 

- Enter a valid username into the search bar.

![Search](https://github.com/irenetony/CS3012-GitHub/raw/master/Screenshots/Search.png) 

- The results should display all the users that username follows along with the number of repositories they have.

![Results](https://github.com/irenetony/CS3012-GitHub/raw/master/Screenshots/Result.png) 

- If you want to enter another username, reload the website and search for another valid github user.

## Built With
- [Go-github](https://github.com/google/go-github): API used to access the github users' data.
- [Gin](https://github.com/gin-gonic/gin): A HTTP web framework written in Go.