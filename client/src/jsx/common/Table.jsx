import React, { useState } from 'react'

import { defaultSort } from '../../utils/sorting'
import { snakeToSentenceCase } from '../../utils/strings'

const Table = ({ rowKey, columns, rows, customRenders, initialSortColumn, initialSortInverse, customSortComparators }) => {
    customRenders ??= {}
    customSortComparators ??= {}

    if (!rows) {
        return <div />
    }

    const [sortColumn, setSortColumn] = useState(initialSortColumn)
    const [sortInverse, setSortInverse] = useState(initialSortInverse ?? false)

    const onHeaderClick = (header_name) => {
        if (sortInverse) {
            setSortColumn(null)
            setSortInverse(false)
        } else if (sortColumn === header_name) {
            setSortInverse(!sortInverse)
        } else {
            setSortColumn(header_name)
            setSortInverse(false)
        }
    }

    let toRender = rows.slice()
    if (sortColumn) {
        toRender.sort((a, b) => (customSortComparators[sortColumn] ?? defaultSort)(
            a[sortColumn],
            b[sortColumn],
        ))
    }
    if (sortInverse) {
        toRender.reverse()
    }

    return (
        <div>
            <table className='styled-table'>
                <thead>
                    <tr>
                        {columns.map(key =>
                            <th
                                key={key}
                                onClick={() => onHeaderClick(key)}
                            >{snakeToSentenceCase(key)}{(key === sortColumn ? (sortInverse ? " ↑" : " ↓") : "")}</th>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {toRender.map((row, _i) =>
                        <tr key={row[rowKey]}>{
                            columns.map(key =>
                                <td key={key}>{
                                    key in customRenders ? customRenders[key](row) : row[key]
                                }</td>
                            )
                        }</tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}

export default Table