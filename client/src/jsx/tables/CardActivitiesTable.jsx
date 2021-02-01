import React from 'react'
import { Link } from 'react-router-dom'

import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
import ExtendableTable from './ExtendableTable'

const CardActivitiesTable = (props) => {
    return <ExtendableTable
        columns={[
            'uuid',
            'transaction_date',
            'post_date',
            'description',
            'category',
            'type',
            'amount'
        ]}
        customRenders={{
            'uuid': (data) =>
                <Link to={'/cards/' + data['card_uuid'] + '/activities/' + data['uuid']}>{data['uuid']}</Link>,
            'transaction_date': (data) => formatAsDate(data['transaction_date']),
            'post_date': (data) => formatAsDate(data['post_date']),
            'amount': (data) => formatAsMoney(data['amount']),
        }}
        initialSortColumn='transaction_date'
        customSortComparators={{
            'transaction_date': dateComparator,
            'post_date': dateComparator,
        }}
        {...props}
    />
}

export default CardActivitiesTable
