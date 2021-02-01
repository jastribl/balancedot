import React, { useState } from 'react'
import { Link } from 'react-router-dom'

import LoaderComponent from '../common/LoaderComponent'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const LinkerPage = () => {
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)

    return (
        <div>
            <h1>Splitwise Expense Linking</h1>
            <LoaderComponent
                path={'/api/splitwise_expenses/unlinked'}
                parentLoading={false}
                setData={setSplitwiseExpenses}
            />
            <SplitwiseExpenseTable
                data={splitwiseExpenses}
                extraColumns={['link']}
                extraCustomRenders={{
                    // todo: try moving linking links to the main splitwise page
                    // todo: style link the rest of the buttons (along with other links)
                    'link': (data) => <Link to={'/linker/' + data['uuid']}>Link</Link>,
                }}
            />
        </div>
    )
}

export default LinkerPage
