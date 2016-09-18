import React from 'react'

export default class Message2 extends React.Component {
  render() {
    return (
      <p>{this.props.data.rank} / language:{this.props.data.language}</p>
    )
  }
}


