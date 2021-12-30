import React from 'react'
import { Snapshot } from '../models/Snapshot'

export default function DateSeparator (props: { snapshot: Snapshot; }) {
  return (
    <div style={{
      marginTop: '30px',
      color: '#eb3448',
      marginBottom: '10px',
      textAlign: 'right',
      fontStyle: 'italic'
    }}> {new Date(props.snapshot.Date).toISOString().split('T')[0]}
    </div>
  )
}
