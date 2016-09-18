import React from 'react'

export default class Options extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
  }

  render() {
    return (
        <option>{this.props.id}</option>
    )
  }
}
