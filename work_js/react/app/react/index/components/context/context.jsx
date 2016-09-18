import React from 'react'

//---------------------------------------------------------
//User Class
//---------------------------------------------------------
export default class User extends React.Component {
  // getChildContext serves as the initializer for our context values
  getChildContext() {
    return {
      favColor: '#f8c483',
      userName: 'James Ipsum'
    }
  }
  
  render() {
    return(
      <div>
        <Usercard />
        <hr/>
      </div>
    )
  }
}

// childContextTypes is defined on the context-provider, giving the context values their corresponding type and passing them down the tree 
User.childContextTypes = {
  favColor: React.PropTypes.string,
  userName: React.PropTypes.string
}

//---------------------------------------------------------
//Usercard Class
//---------------------------------------------------------
class Usercard extends React.Component {
  // Note that the Usercard component makes no use of context
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
  // We can make use of these context values...
  render() {
    return(
      <h2>{this.context.userName}</h2>
    )
  }
}

// By defining corresponding contextTypes on child components down the tree that wish to access them
UserInfo.contextTypes = {
  userName: React.PropTypes.string
}

//---------------------------------------------------------
//UserIcon Class
//---------------------------------------------------------
class UserIcon extends React.Component {
  // Same as above
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

// But accessing favColor instead.
UserIcon.contextTypes = {
  favColor: React.PropTypes.string
}
