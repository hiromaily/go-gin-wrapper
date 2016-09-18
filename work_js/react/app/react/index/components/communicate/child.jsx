import React from 'react'

export default class Child extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
      input: ''
    }
  }

 render() {
    return (
      <h2 style={{color: this.props.color}}>
        {this.props.currentInput}
      </h2>
    )
  }
}
