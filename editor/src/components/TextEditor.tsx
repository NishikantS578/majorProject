import React from 'react'

const TextEditor = () => {
  return (
    <div>
      <div className="menubar">
        <button id="open">Open</button>
        <button id="save">Save</button>
        <button id="compileAndRun">Compile and Run</button>
      </div>
      <textarea className="editorTab"></textarea>
      <div id="outputWindow">
        <h3>OUTPUT</h3>
        <textarea id="output" disabled></textarea>
      </div>
    </div>
  )
}

export default TextEditor