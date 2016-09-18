import React from 'react'

export default class Increment extends React.Component {
  render() {
    return(
      <button onClick={this.props.increment}>Higher</button>
    )
  }
}
