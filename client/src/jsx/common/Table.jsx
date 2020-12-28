import React from 'react'

const Table = ({ rowKey, columns, rows, customRenders }) => {
    if (customRenders === undefined) {
        customRenders = {}
    }
    if (!rows) {
        return <div />
    }
    return (
        <div>
            <table className='styled-table'>
                <thead>
                    <tr>
                        {Object.keys(columns).map((key) => <th key={key}>{columns[key]}</th>)}
                    </tr>
                </thead>
                <tbody>
                    {rows.map((row, i) =>
                        <tr key={row[rowKey]}>{
                            Object.keys(columns).map((key) =>
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