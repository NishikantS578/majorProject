import { useState } from 'react'
import Blocks from './blocks'

const save = (editor_content_state: any) => {
  console.log(editor_content_state)
  window.ipcRenderer.invoke("saveFile", editor_content_state)
}

const open = () => {
  window.ipcRenderer.invoke("openFile")
}

const compileRun = () => {
  window.ipcRenderer.invoke("compileAndRun")
}

function render_flow_chart(node: any): any {
  if (node instanceof Blocks.StartBlock) {
    return <>
      {node.Block()}
      {render_flow_chart(node.child_block)}
    </>
  } else if (node instanceof Blocks.ProcessingBlock) {
    return <>
      {node.Block()}
      {render_flow_chart(node.child_block)}
    </>
  } else if (node instanceof Blocks.IfBlock) {
    return <>
      {node.Block()}
      {render_flow_chart(node.true_child_block)}
      {render_flow_chart(node.false_child_block)}
    </>
  } else if (node instanceof Blocks.EndBlock) {
    return <>
      {node.Block()}
    </>
  }
}

function App() {
  const [editor_content, set_editor_content_state] = useState("")
  const [output_content, set_output_content] = useState("")
  const [window_type, set_window_type] = useState("text_editor")

  window.ipcRenderer.on("openedFile", (_, data) => {
    set_editor_content_state(data)
  })

  window.ipcRenderer.on("compiledAndRun", (_, data) => {
    set_output_content(data)
  })

  let start_block = new Blocks.StartBlock(null)
  let processing_block = new Blocks.ProcessingBlock(null, null)
  start_block.child_block = processing_block
  let end_block = new Blocks.EndBlock()
  processing_block.child_block = end_block

  return (
    <>
      <div className="menubar">
        <button onMouseDown={open}>Open</button>
        <button onMouseDown={() => save(editor_content)}>
          Save
        </button>
        <button onMouseDown={compileRun}>
          Compile and Run
        </button>
        <button onClick={() => {
          if (window_type == "text_editor") {
            set_window_type("visual_editor")
          } else {
            set_window_type("text_editor")
          }
        }}>
          {window_type}
        </button>
      </div>

      <div className='editor min-w-fit w-full flex-auto'>
        {
          window_type == "text_editor" ?
            (
              <textarea className="editorTab min-w-fit text-nowrap w-full h-full"
                onChange={
                  (e) => set_editor_content_state(e.target.value)}
                value={editor_content}>
              </textarea>
            )
            :
            (
              <div className='visualEditor w-full h-full'>
                {
                  render_flow_chart(start_block)
                }
              </div>
            )
        }
      </div>

      <div className='outputWindow'>
        <h3>OUTPUT</h3>
        <textarea className='output'
          disabled value={output_content}>
        </textarea>
      </div>
    </>
  )
}

export default App