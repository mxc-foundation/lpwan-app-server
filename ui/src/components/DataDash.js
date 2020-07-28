import React from 'react'
import styled from 'styled-components'
import i18n, { packageNS } from '../i18n';
import phone from '../assets/images/appimage@2x.png';
import dataDash from '../assets/images/dataDash.png'
import appStore from '../assets/images/Appstore@2x.png';
import googlePlay from '../assets/images/google-play.png';
import testflight from '../assets/images/testflight.png';
import downloadAPK from '../assets/images/download_apk.png';

import { Card } from 'reactstrap';



const DataDash = () => {
    return (
        <Card>
            <Flexbox>
              <Images>
            <Image src={dataDash} width="50"/>
            </Images>
            <Title>{i18n.t(`${packageNS}:m2m_redirect.title`)}</Title>
            <Subtitle>{i18n.t(`${packageNS}:m2m_redirect.subtitle`)} </Subtitle>
            <Description>{i18n.t(`${packageNS}:m2m_redirect.description`)}</Description>
            <FlexRow>
            <Buttons>
            <Button href="https://apps.apple.com/app/mxc-datadash/id1509218470"><img src={appStore} width="135"/></Button>
            <BorderButton href="https://testflight.apple.com/join/NkXHEpf4"><img src={testflight} height="39.844"/>  Install with TestFlight</BorderButton>
            <Button href="https://play.google.com/store/apps/details?id=com.mxc.smartcity"><img src={googlePlay} width="135"/></Button>
            <Button href="https://datadash.oss-accelerate.aliyuncs.com/app-prod-release.apk"><img src={downloadAPK} width="135"/></Button>
            </Buttons>
            </FlexRow>
            <Images>
            <Image src={phone} width="600"/>
            </Images>
            </Flexbox>
            </Card>
    )
}

const Flexbox = styled.div`
display: flex;
flex-direction: column;
padding: 4vw 4vw 0 4vw;
`

const FlexRow = styled.div`
display: flex;
flex-direction: row;
justify-content: space-around;
align-items: space-around;
`

const Title = styled.h2`
 display: flex;
 justify-content: center;
 align-items: center;
 color: #1C1478;
 font-size: 1.25rem;
`

const Subtitle = styled.p`
justify-content: center;
align-items: center;
text-align: center;
font-size: 18px;
font-weight: bold;
`

const Description = styled.p`
justify-content: center;
align-items: center;
text-align: center;
font-size: 18px;
margin: 4vh 0 8vh 0;
`
const Buttons = styled.div`
justify-content: space-around;
align-items: center;
margin: 0 0 10vh 0;
`

const Button = styled.a`
align-items: space-around;
margin: 3vw;
color: black;

`

const BorderButton = styled.a`
align-items: space-around;
margin: 3vw;
color: black;
border: 2px solid black;
border-radius: 11px;
padding: 11px 5px 12px 0px;
z-index: 5;
`


const Images = styled.div`
 display: flex;
 justify-content: center;
 align-items: center;
 color: #1C1478;
 font-size: 1.25rem;
`

const Image = styled.img`
display: flex;
justify-content: center;
align-items: center;
text-align: center;
`
export default DataDash
