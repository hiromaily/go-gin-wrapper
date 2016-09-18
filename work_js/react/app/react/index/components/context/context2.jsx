import React from 'react'

//---------------------------------------------------------
//User Class
//---------------------------------------------------------
export default class User2 extends React.Component {
  constructor(props){
    super(props)
    
    // Move the currentColor to top-level state
    this.state={
      currentColor: '#f8c483'
    }
  }
  
  // Hook context favColor to this.state.currentColor
  getChildContext() {
    return {
      favColor: this.state.currentColor,
      userName: 'James Ipsum'
    }
  }
  
  // Swap between colors on click
  swapColor () {
    if(this.state.currentColor === '#f8c483'){
      this.setState({
        currentColor: '#76daff'
      })
    } else {
      this.setState({
        currentColor: '#f8c483'
      })
    }
  }
  
  render() {
    return(
      <div>
        <Usercard />
        <button onClick={this.swapColor.bind(this)}>Swap Color</button>
        <hr/>
      </div>
    )
  }
}

User2.childContextTypes = {
  favColor: React.PropTypes.string,
  userName: React.PropTypes.string
}

//---------------------------------------------------------
//Usercard Class
//---------------------------------------------------------
class Usercard extends React.Component {
  // Usercard again makes no use of context
  render() {
    return(
      <div className='usercard'>
        <UserIcon />
        <UserInfo />
      </div>
    )
  }
}

//---------------------------------------------------------
//UserInfo Class
//---------------------------------------------------------
class UserInfo extends React.Component {
  render() {
    return(
      <h2>{this.context.userName}</h2>
    )
  }
}

UserInfo.contextTypes = {
  userName: React.PropTypes.string
}

//---------------------------------------------------------
//UserIcon Class
//---------------------------------------------------------
class UserIcon extends React.Component {
  // The context value for favColor now correctly corresponds to the currentColor value in the User component state
  render() {
    return(
      <div 
        className='circle' 
        style={{
          backgroundColor: this.context.favColor
        }}></div>
    )
  }
}

UserIcon.contextTypes = {
  favColor: React.PropTypes.string
}
