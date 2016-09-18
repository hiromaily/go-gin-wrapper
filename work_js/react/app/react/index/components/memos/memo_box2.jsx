import React from 'react'
import Message2   from './message2.jsx'

export default class MemoBox2 extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        data: {
          rank: 0,
          language: "PHP"
        }
    }
  }

  render() {
    //child node
    let commentNodes = this.props.data.map(function (comment) {
      return (
        //<Message2 data={this.props.data} />
        <Message2 key={comment.id} data={comment} />
      )
    })
    //parent node
    return (
      <div className="commentList">
        <h4>I like ...</h4>
        {commentNodes}
      </div>
    )
  }
}
