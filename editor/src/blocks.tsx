class StartBlock {
    child_block: any
    constructor(child_block: any) {
        this.child_block = child_block
    }

    Block() {
        return <div style={{ display: "flex", alignItems: "center", justifyContent: "center" }} className="flex-col">
            <div style={{ padding: "10px 20px", border: "1px solid white", borderRadius: "100px", display: "flex", alignItems: "center", justifyContent: "center" }}>start block</div>
            <div className="w-[2px] bg-white h-5"></div>
        </div >
    }
}

class EndBlock {
    Block() {
        return <div style={{ display: "flex", alignItems: "center", justifyContent: "center" }} className="flex-col">
            <div className="w-[2px] bg-white h-5"></div>
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
        return <div style={{ display: "flex", alignItems: "center", justifyContent: "center" }} className="flex-col">
            <div className="w-[2px] bg-white h-5"></div>
            <div style={{ border: "1px solid white", padding: "10px 20px" }}>
                {this.statements}
            </div>
            <div className="w-[2px] bg-white h-5"></div>
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
        return <div className="relative h-28 flex items-center justify-center ">
            {this.condition}
            <div className="border border-white rotate absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rotate-45 aspect-square h-[100px] -skew-x-[14deg] -skew-y-[8deg]"></div>
        </div>
    }
}

function ast_to_blocks(ast: any): any {
    return null
}

export default { StartBlock, EndBlock, ProcessingBlock, IfBlock, ast_to_blocks }