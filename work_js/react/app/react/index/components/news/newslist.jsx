import React from 'react'

export default class NewsList extends React.Component {
  render() {
    return (
      <div>
        <h5>[{this.props.date}] {this.props.title}</h5>
        <p> {this.props.children}</p>
        <p>written by {this.props.author}</p>
      </div>
    )
  }
}
