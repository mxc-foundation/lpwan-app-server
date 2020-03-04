import React from "react";
import { Row, Col, Button } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';


const Phrase = ({ srNo, phrase }) => {
    return <React.Fragment>
        <span className="phrase-item">
            <span className="no">{srNo}.</span><span className="word">{phrase}</span>
        </span>
    </React.Fragment>
}

const MneMonicPhrase = ({ title, phrase, showSkip, next, close }) => {
    let list = [...(phrase || [])];
    const chunks = [];
    let chunkCount = 3;
    while (list.length) {
        const chunkSize = Math.ceil(list.length / chunkCount--);
        const chunk = list.slice(0, chunkSize);
        chunks.push(chunk);
        list = list.slice(chunkSize);
    }

    return <React.Fragment>
        <Row className="text-center">
            <Col className="mb-0">
                <h4>{title}</h4>

                <Row className="text-left mt-3">
                    {chunks.length ? <Col>
                        <ul className="list-unstyled">
                            {chunks[0].map((word, idx) => {
                                return <li key={idx}><Phrase srNo={idx + 1} phrase={word} /></li>
                            })}
                        </ul>
                    </Col> : null}

                    {chunks.length > 1 ? <Col>
                        <ul className="list-unstyled">
                            {chunks[1].map((word, idx) => {
                                return <li key={idx}><Phrase srNo={chunks[0].length + idx + 1} phrase={word} /></li>
                            })}
                        </ul>
                    </Col> : null}

                    {chunks.length > 2 ? <Col>
                        <ul className="list-unstyled">
                            {chunks[2].map((word, idx) => {
                                return <li key={idx}><Phrase srNo={chunks[0].length + chunks[1].length + idx + 1} phrase={word} /></li>
                            })}
                        </ul>
                    </Col> : null}
                </Row>

                <Button color="primary" className="btn-block mt-2" onClick={next}>{i18n.t(`${packageNS}:menu.menmonic_phrase.write_button`)}</Button>
                {showSkip ? <Button color="link" className="btn-block" onClick={close}>{i18n.t(`${packageNS}:menu.menmonic_phrase.skip_button`)}</Button> : null}
            </Col>
        </Row>
    </React.Fragment>
}

export default MneMonicPhrase;
