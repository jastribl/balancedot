import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { getWithHandling } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const LinkerPage = () => {
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)
    const [pageLoading, setPageLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    useEffect(() => {
        getWithHandling(
            '/api/splitwise_expenses/unlinked',
            setSplitwiseExpenses,
            setErrorMessage,
            setPageLoading
        )
    }, [
        setPageLoading,
        setSplitwiseExpenses,
        setErrorMessage,
    ])

    return (
        <div>
            <Spinner visible={pageLoading} />
            <h1>Splitwise Expense Linking</h1>
            <ErrorRow message={errorMessage} />
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
