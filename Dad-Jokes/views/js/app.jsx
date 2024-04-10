class App extends React.Component{
    render(){
        if(this.loggedIn){
            return (<loggedIn/>)
        } else{
            return (<Home/>)
        }
    }
}

//HOME COMPONENT
class Home extends React.Component{
    render(){
        return (
            <div className="container">
                <div className="cols-xs-8 col-xs-offset-2 jumbotron text-center">
                    <h1>Jokeish</h1>
                    <p>A load of Dad Jokes</p>
                    <p>Sign in to get access</p>
                    <a onClick={this.authenticate} className="btn btn-primary btn-lg btn-login btn-block">SignIn</a>
                </div>
            </div>
        )
    }
}

//LOGGEDIN COMPONENT
class loggedIn extends React.Component{
    constructor(props){
        super(props)
        this.state={
            jokes: []
        }
    }

    render(){
        return (
            <div className="container">
                <div className="col-lg-12">
                    <br />
                    <span className="pull-right"><a onClick={this.logout}>Logout</a></span>
                    <h2>Jokeish</h2>
                    <p>Lets feed you with some funny jokes</p>
                    <div className="row">
                        {this.state.jokes.map(function(joke,i){
                            return (<Joke Key={i} joke={joke}/>)
                        })}
                    </div>

                </div>
            </div>
        )
    }
}

//JOKE COMPONENT
class Joke extends React.Component{
    constructor(props){
        super(props)
        this.state = {
            liked: ""
        }
        this.like = this.like.bind(this);
    }

    like(){
        //TODO
    }

    render(){
        return(
            <div className="col-xs-4">
                <div className="panel panel-default">
                    <div className="panel-heading">#{this.props.joke.id}<span className="pull-right">{this.state.liked}</span></div>
                    <div className="panel-body">{this.props.joke.joke}</div>
                    <div className="panel-folder">{this.props.joke.likes} Like &nbsp; <a onClick={this.like} className="btn btn-default"><span className="glyphicon glyphicon-thumbs-up"></span></a></div>
                </div>
            </div>
        )
    }
}

ReactDOM.render(<App />, document.getElementById('app'));

