import classNames from "classnames";
import React, { useState } from "react";
import { Button, Col, Row } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';



const Phrase = ({ phrase, isSelected, select }) => {
    return <React.Fragment>
        <Button color={isSelected ? "primary" : "secondary"} outline className={classNames("btn-rounded", "m-1", { "bg-white": isSelected, "text-primary": isSelected })} onClick={select}>
            {phrase}
        </Button>
    </React.Fragment>
}

const shuffleArray = (array) => {
    for (let i = array.length - 1; i > 0; i--) {
        let j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
}

const MneMonicPhraseConfirm = ({ title, phrase, next, back, skip, showBackButton = true, showSkipButton = false, titleClass = "" }) => {

    const [selectedPhrase, setSelectedPhrase] = useState([]);

    const onSelect = (phrase, isSelected) => {
        let phrases = [...selectedPhrase];
        if (isSelected) {
            phrases = phrases.filter(p => p !== phrase);
        } else {
            phrases.push(phrase);
        }
        setSelectedPhrase(phrases);
    }

    const randomizedPhraseList = shuffleArray([...(phrase || [])]);

    return <React.Fragment>
        <Row className="text-center">
            <Col className="mb-0">
                <h5 className={titleClass}>{title}</h5>

                <Row className="mt-3 text-left">
                    <Col className="mb-0">
                        <div className="bg-light p-3">
                            {selectedPhrase.map((word, idx) => (
                                <Phrase key={idx} phrase={word} isSelected={true} select={() => { onSelect(word, true) }} />
                            ))}
                        </div>
                    </Col>
                </Row>

                <Row className="mt-3 text-left">
                    <Col className="mb-0">
                        <div>
                            {randomizedPhraseList.map((word, idx) => {
                                return selectedPhrase.indexOf(word) === -1 ? <Phrase key={idx} phrase={word} isSelected={false} select={() => { onSelect(word, false) }} />: null
                            })}
                        </div>
                    </Col>
                </Row>

                <Row className="mt-2 text-left">
                    <Col className="mb-0">
                        <Button color="primary" className="btn-block" onClick={() => next(selectedPhrase)}
                            disabled={!selectedPhrase.length}>{i18n.t(`${packageNS}:menu.menmonic_phrase.confirm_button`)}</Button>
                    </Col>
                    {showBackButton ? <Col className="mb-0">
                        <Button color="primary" outline className="btn-block" onClick={back}>{i18n.t(`${packageNS}:menu.menmonic_phrase.back_button`)}</Button>
                    </Col> : null}
                </Row>

                {showSkipButton ? <Row className="mt-2 text-left">
                    <Col className="mb-0">
                        <Button color="link" className="btn-block" onClick={skip}>{i18n.t(`${packageNS}:menu.menmonic_phrase.skip_button`)}</Button>
                    </Col>
                </Row> : null}

            </Col>
        </Row>
    </React.Fragment>
}

export default MneMonicPhraseConfirm;
