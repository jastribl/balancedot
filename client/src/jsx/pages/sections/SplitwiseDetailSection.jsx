import React, { useState } from 'react'

import LoaderComponent from '../../common/LoaderComponent'

const SplitwiseDetailSection = ({ splitwiseExpenseID }) => {
    const [raw, setRaw] = useState(null)

    let dataSection = null
    if (raw) {
        dataSection = <div>
            <pre>{
                JSON.stringify(JSON.parse(raw)['expense'], null, 4)
            }</pre>
        </div>
    }

    return <div>
        <LoaderComponent
            path={`/api/splitwise_expenses/${splitwiseExpenseID}/raw`}
            parentLoading={false}
            setData={setRaw}
        />
        {dataSection}
    </div>
}

export default SplitwiseDetailSection
