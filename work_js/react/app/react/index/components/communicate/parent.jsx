import React from 'react'
import Child   from './child.jsx'

export default class Parent extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
      input: ''
    }
  }

  handleChange(e) {
    this.setState({
      input: e.target.value
    })
  }
  //parent to child
  render() {
    return (
      <div>
        <h4>Parent to Child</h4>
        <input 
          onChange={this.handleChange.bind(this)} 
          value={this.state.input} 
          placeholder="Type away" />
        <Child 
          currentInput={this.state.input} 
          color="#76daff" />
      </div>
    )
  }
}
