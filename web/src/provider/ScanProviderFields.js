// Copyright 2026 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from "react";
import {Button, Col, Input, Row} from "antd";
import * as Setting from "../Setting";
import i18next from "i18next";
import ScanTable from "../table/ScanTable";

const {TextArea} = Input;

export function renderScanProviderFields(provider, updateProviderField, options = {}) {
  const canScan = options.mode !== "add";
  if (provider.type === "Security Scan") {
    return (
      <React.Fragment>
        <Row style={{marginTop: "20px"}}>
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("provider:Online list")}:
          </Col>
          <Col span={22}>
            <Input value={provider.endpoint} onChange={e => updateProviderField("endpoint", e.target.value)} />
          </Col>
        </Row>
        {provider.subType === "Url" ? (
          <Row style={{marginTop: "20px"}}>
            <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
              {i18next.t("general:URL")}:
            </Col>
            <Col span={22}>
              <TextArea
                autoSize={{minRows: 3, maxRows: 10}}
                value={provider.content}
                placeholder="https://example.com\nhttps://another.example.com"
                onChange={e => updateProviderField("content", e.target.value)}
              />
            </Col>
          </Row>
        ) : null}
        <Row style={{marginTop: "20px"}}>
          <Col span={22} offset={(Setting.isMobile()) ? 0 : 2}>
            <Button
              type="primary"
              loading={options.scanLoading}
              disabled={!canScan}
              onClick={() => options.onScan(provider.subType === "Url" ? provider.content : "")}
            >
              {i18next.t("general:Scan")}
            </Button>
          </Col>
        </Row>
        <ScanTable provider={provider} options={{...options, subType: provider.subType, owner: provider.owner}} />
      </React.Fragment>
    );
  }

  return null;
}
