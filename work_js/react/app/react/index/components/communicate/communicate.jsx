import React from 'react'
import Parent   from './parent.jsx'

export default class Communicate extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        count: 0
    }
  }

  render() {
    return (
      <div>
        <Parent />
        <hr/>
      </div>
    )
  }
}
