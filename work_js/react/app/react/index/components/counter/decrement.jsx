import React from 'react'

export default class Decrement extends React.Component {
  // And again for this.props.decrement
  render() {
    return(
      <button onClick={this.props.decrement}>Lower</button>
    )
  }
}
