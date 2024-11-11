class StartBlock {
    child_block: any
    constructor(child_block: any) {
        this.child_block = child_block
    }

    Block() {
        return <div style={{ display: "flex", alignItems: "center", justifyContent: "center" }}>
            <div style={{ padding: "10px 20px", border: "1px solid white", borderRadius: "100px", display: "flex", alignItems: "center", justifyContent: "center" }}>start block</div>
        </div >
    }
}

class EndBlock {
    Block() {
        return <div style={{ display: "flex", alignItems: "center", justifyContent: "center" }}>
            <div style={{ padding: "10px 20px", border: "1px solid white", borderRadius: "100px", display: "flex", alignItems: "center", justifyContent: "center" }}>end block</div>
        </div >
    }
}

class ProcessingBlock {
    statements: string
    child_block: any
    constructor(statements: any, child_block: any) {
        this.statements = statements
        this.child_block = child_block
    }

    Block() {
        return <div style={{ display: "flex", alignItems: "center", justifyContent: "center" }}>
            <div style={{border: "1px solid white", padding: "10px 20px"}}>
                processing block
            </div>
        </div>
    }
}

class IfBlock {
    condition: string
    true_child_block: any
    false_child_block: any
    constructor(condition: any, true_child_block: any, false_child_block: any) {
        this.condition = condition
        this.true_child_block = true_child_block
        this.false_child_block = false_child_block
    }

    Block() {
        return <>if block</>
    }
}

function ast_to_blocks(ast: any): any{
    return null
}

export default { StartBlock, EndBlock, ProcessingBlock, IfBlock, ast_to_blocks }