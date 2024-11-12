import { useState } from 'react'
import Blocks from './blocks'
// import parser from './parser'

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
      <div className="">
        {render_flow_chart(node.true_child_block)}
        {render_flow_chart(node.false_child_block)}
      </div>
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
  let processing_block2 = new Blocks.ProcessingBlock(null, null)
  start_block.child_block = processing_block
  let end_block = new Blocks.EndBlock()
  processing_block.child_block = processing_block2
  processing_block2.child_block = end_block
  processing_block.statements = "let a = 11"
  processing_block2.statements = "let b = 12"
  let if_block = new Blocks.IfBlock('11 == 11', null, null)
  processing_block2.child_block = if_block
  let processing_block3 = new Blocks.ProcessingBlock(null, null)
  processing_block3.statements = "let c = 13"
  if_block.true_child_block = processing_block3
  processing_block3.child_block = end_block
  let flow_chart = start_block

  // let ast = parser(editor_content)
  // let flow_chart = Blocks.ast_to_blocks(ast)

  return (
    <>
      <div className="menubar flex justify-between">
        <div className='flex menubar'>
          <button onMouseDown={open}>Open</button>
          <button onMouseDown={() => save(editor_content)}>
            Save
          </button>
          <button onMouseDown={compileRun}>
            Compile and Run
          </button>
        </div>

        <div className="flex menubar">
          <button onClick={() => { set_window_type("text_editor") }} className={window_type == "text_editor" ? "active-window" : ""}>
            Text Editor
          </button>
          <button onClick={() => { set_window_type("visual_editor") }} className={window_type == "visual_editor" ? "active-window" : ""}>
            Visual Editor
          </button>
        </div>
      </div>

      <div className='editor min-w-fit w-full flex-auto'>
        {
          window_type == "text_editor" ?
            (
              <textarea className="editorTab min-w-fit text-nowrap w-full h-full text-lg p-4"
                onChange={
                  (e) => set_editor_content_state(e.target.value)}
                value={editor_content}>
              </textarea>
            )
            :
            (
              <div className='visualEditor w-full h-full flex flex-col'>
                {
                  render_flow_chart(flow_chart)
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