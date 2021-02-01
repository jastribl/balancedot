import React from 'react'
import { Link } from 'react-router-dom'

import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
import ExtendableTable from './ExtendableTable'

const SplitwiseExpenseTable = (props) => {
    return <ExtendableTable
        columns={[
            'uuid',
            'splitwise_id',
            'description',
            'details',
            'creation_method',
            'amount',
            'amount_paid',
            'date',
            'category',
        ]}
        customRenders={{
            'uuid': (data) => <Link to={'/splitwise_expenses/' + data['uuid']}>{data['uuid']}</Link>,
            'details': (data) => data['details'].trim(),
            'date': (data) => formatAsDate(data['date']),
            'amount': (data) => formatAsMoney(data['amount'], data['currency_code']),
            'amount_paid': (data) => formatAsMoney(data['amount_paid'], data['currency_code']),
        }}
        initialSortColumn='date'
        initialSortInverse={true}
        customSortComparators={{
            'date': dateComparator
        }}
        {...props}
    />
}

export default SplitwiseExpenseTable
