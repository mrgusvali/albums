import logo from './logo.svg';
import './App.css';
import React, { Component } from 'react';

class Form extends React.Component {
	constructor(props) {
		super(props)
		
		this.state = {t: "", a: ""}
	}
	
	handleTitle = (ev) => {
		console.log("#handleTitle", ev.target.value)
		this.setState({t: ev.target.value})
		this.props.onChange(this.state)	
	}
	handleArtist = (ev) => {
		console.log("#handleArtist", ev.target.value)
		this.setState({a: ev.target.value})
		this.props.onChange(this.state)	
	}	
  render() { return (
			<div className="form">
			  <label>Title: <input type="text" onChange={this.handleTitle} /></label>
			  <p/>
			  <label>Artist: <input type="text" onChange={this.handleArtist} /></label>
			</div>
		);
	}
}

/* listing preview of messages  */
function MessagesList(props) {
    const { error, isLoaded, albums } = props.loadedState;
	console.log("loaded state", error, isLoaded, albums)
    if (error) {
      return <div>Albums List Error: {error.message}. Is the web service running?</div>;
    } else if (!isLoaded) {
      return <div>Loading...</div>;
    } else if (albums.length == 0) {
      return <div>An empty list...</div>;    
    } else {
      return (
        <div id="AlbumsList">
         {albums.map(a=><label key={a.id}>{a.title}:{a.artist}<br/></label>)}
        </div>
      );
    }
}

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
	  filterCriteria: {t:"", a:"", offset:0},
	  loadedState: {
        error: null,
        isLoaded: false,
        albums: []
      }
    };
  }		
	
  componentDidMount() {
	this.fetch(this.state.filterCriteria)
  }	
	
  fetch(criteria) {
	console.log("#fetch", criteria)
	let opts = {
		method: 'GET'
	}
	var url = 'http://localhost:8080/albums?' + new URLSearchParams(criteria).toString()
	console.log("url ", url)
    fetch(url, opts)
      .then(res => res.json())
      .then(
        (result) => {
          console.log(result.length + " messages.")
          this.setState({loadedState: {
            isLoaded: true,
            albums: result
          }});
        },
        (error) => {
          console.log("error: ", error) 	
          this.setState({loadedState: {
            isLoaded: true,
            error
          }});
        }
      )
  }	
  
  handleFormChange = (criteria)=>{this.fetch(criteria)}
	
  render() { return (
	    <div className="App">
	      <header className="App-header">
	        <img src={logo} className="App-logo" alt="logo" />
	        <p>
	          Edit <code>src/App.js</code> and save to reload.
	        </p>
	      </header>
	      
	      <Form onChange={this.handleFormChange} />
	      <MessagesList loadedState={this.state.loadedState} />
	    </div>
	  );
	}
}

export default App;
