import React from 'react'

export default class Message extends React.Component {
  render() {
    return (
      <div>  
        <h5>[about your schedule]</h5>
        <p>day:{this.props.day} / title:{this.props.title}</p>
        <p>{this.props.children}</p>
      </div>
    )
  }
}
