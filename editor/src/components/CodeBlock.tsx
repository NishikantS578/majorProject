
function CodeBlock(props: {
    graph?: Node,
    arrow_start_x: number, arrow_start_y: number
}) {
    if (props.graph === undefined) {
        return
    }
    return <svg>
    </svg>
}

export default CodeBlock