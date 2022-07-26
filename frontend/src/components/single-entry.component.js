import React, {useState, useEffect} from 'react';

import 'bootstrap/dist/css/bootstrap.css';

import { Card, Button, Row, Col } from 'react-bootstrap';

const Entry = ({entryData, setChangeIngredient, deleteSingleEntry, setChangeEntry}) => {
    return (
        <Card>
             <Row>
                <Col>Dish:{entryData !== undefined && entryData.dish}</Col>
             </Row>
        </Card>
    )
}