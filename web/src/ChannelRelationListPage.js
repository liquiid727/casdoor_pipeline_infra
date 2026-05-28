import React from "react";
import {Button, Card, Col, Input, InputNumber, Row, Select, Switch, Table} from "antd";
import * as ChannelRelationBackend from "./backend/ChannelRelationBackend";
import * as OrganizationBackend from "./backend/OrganizationBackend";
import * as Setting from "./Setting";
import PopconfirmModal from "./common/modal/PopconfirmModal";
import i18next from "i18next";

const {Option} = Select;

class ChannelRelationListPage extends React.Component {
  constructor(props) {
    super(props);
    const randomName = Setting.getRandomName();
    this.state = {
      organizationName: props.match.params.organizationName,
      organizations: [],
      relations: [],
      loading: true,
      editingRelation: {
        owner: "admin",
        name: `channel_relation_${randomName}`,
        parentOrganization: props.match.params.organizationName,
        channelOrganization: "",
        channelMode: "shared_account",
        settlementRule: "",
        settlementRatio: 0,
        settlementAccount: "",
        status: "active",
      },
    };
  }

  UNSAFE_componentWillMount() {
    this.getOrganizations();
    this.getRelations();
  }

  getOrganizations() {
    OrganizationBackend.getOrganizations("admin")
      .then((res) => {
        if (res.status === "ok") {
          this.setState({organizations: res.data || []});
        } else {
          Setting.showMessage("error", res.msg);
        }
      });
  }

  getRelations() {
    this.setState({loading: true});
    ChannelRelationBackend.getChannelRelations("admin", this.state.organizationName)
      .then((res) => {
        this.setState({loading: false});
        if (res.status === "ok") {
          this.setState({relations: res.data || []});
        } else {
          Setting.showMessage("error", res.msg);
        }
      });
  }

  updateRelationField(key, value) {
    const editingRelation = {...this.state.editingRelation};
    editingRelation[key] = value;
    this.setState({editingRelation});
  }

  editRelation(relation) {
    this.setState({
      editingRelation: {
        owner: relation.owner,
        name: relation.name,
        parentOrganization: relation.parentOrganization,
        channelOrganization: relation.channelOrganization,
        channelMode: relation.channelMode || "shared_account",
        settlementRule: relation.settlementRule || "",
        settlementRatio: relation.settlementRatio || 0,
        settlementAccount: relation.settlementAccount || "",
        status: relation.status || "active",
      },
    });
  }

  newRelation() {
    const randomName = Setting.getRandomName();
    this.setState({
      editingRelation: {
        owner: "admin",
        name: `channel_relation_${randomName}`,
        parentOrganization: this.state.organizationName,
        channelOrganization: "",
        channelMode: "shared_account",
        settlementRule: "",
        settlementRatio: 0,
        settlementAccount: "",
        status: "active",
      },
    });
  }

  submitRelation() {
    const relation = Setting.deepCopy(this.state.editingRelation);
    const exists = this.state.relations.some(item => item.name === relation.name);
    const action = exists
      ? ChannelRelationBackend.updateChannelRelation("admin", relation.name, relation)
      : ChannelRelationBackend.addChannelRelation(relation);

    action.then((res) => {
      if (res.status === "ok") {
        Setting.showMessage("success", i18next.t("general:Successfully saved"));
        this.getRelations();
        if (!exists) {
          this.newRelation();
        }
      } else {
        Setting.showMessage("error", res.msg);
      }
    });
  }

  deleteRelation(relation) {
    ChannelRelationBackend.deleteChannelRelation(relation)
      .then((res) => {
        if (res.status === "ok") {
          Setting.showMessage("success", i18next.t("general:Successfully deleted"));
          this.getRelations();
          this.newRelation();
        } else {
          Setting.showMessage("error", res.msg);
        }
      });
  }

  renderEditor() {
    const relation = this.state.editingRelation;
    const channelOptions = this.state.organizations.filter(org => org.name !== this.state.organizationName);

    return (
      <Card size="small" title={i18next.t("general:Edit")}>
        <Row style={{marginTop: "10px"}}>
          <Col span={4}>{i18next.t("general:Name")}:</Col>
          <Col span={20}>
            <Input value={relation.name} onChange={e => this.updateRelationField("name", e.target.value)} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={4}>{i18next.t("general:Organization")}:</Col>
          <Col span={20}>
            <Input value={relation.parentOrganization} disabled />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={4}>Channel:</Col>
          <Col span={20}>
            <Select virtual={false} style={{width: "100%"}} value={relation.channelOrganization} onChange={value => this.updateRelationField("channelOrganization", value)}>
              {channelOptions.map(org => <Option key={org.name} value={org.name}>{org.displayName || org.name}</Option>)}
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={4}>{i18next.t("general:Type")}:</Col>
          <Col span={20}>
            <Select virtual={false} style={{width: "100%"}} value={relation.channelMode} onChange={value => this.updateRelationField("channelMode", value)}>
              <Option value="shared_account">shared_account</Option>
              <Option value="sub_tenant">sub_tenant</Option>
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={4}>{i18next.t("general:Rule")}:</Col>
          <Col span={20}>
            <Input value={relation.settlementRule} onChange={e => this.updateRelationField("settlementRule", e.target.value)} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={4}>Settlement ratio:</Col>
          <Col span={20}>
            <InputNumber style={{width: "100%"}} value={relation.settlementRatio} onChange={value => this.updateRelationField("settlementRatio", value ?? 0)} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={4}>{i18next.t("general:Account")}:</Col>
          <Col span={20}>
            <Input value={relation.settlementAccount} onChange={e => this.updateRelationField("settlementAccount", e.target.value)} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={4}>{i18next.t("general:Status")}:</Col>
          <Col span={20}>
            <Select virtual={false} style={{width: "100%"}} value={relation.status} onChange={value => this.updateRelationField("status", value)}>
              <Option value="active">active</Option>
              <Option value="disabled">disabled</Option>
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={24}>
            <Button type="primary" onClick={() => this.submitRelation()}>{i18next.t("general:Save")}</Button>
            <Button style={{marginLeft: "12px"}} onClick={() => this.newRelation()}>{i18next.t("general:Add")}</Button>
          </Col>
        </Row>
      </Card>
    );
  }

  renderTable() {
    const columns = [
      {title: i18next.t("general:Name"), dataIndex: "name", key: "name"},
      {title: "Channel", dataIndex: "channelOrganization", key: "channelOrganization"},
      {title: i18next.t("general:Type"), dataIndex: "channelMode", key: "channelMode"},
      {title: i18next.t("general:Rule"), dataIndex: "settlementRule", key: "settlementRule"},
      {title: "Settlement ratio", dataIndex: "settlementRatio", key: "settlementRatio"},
      {title: i18next.t("general:Account"), dataIndex: "settlementAccount", key: "settlementAccount"},
      {
        title: i18next.t("general:Status"),
        dataIndex: "status",
        key: "status",
        render: value => <Switch checked={value !== "disabled"} disabled />,
      },
      {
        title: i18next.t("general:Action"),
        key: "action",
        render: (_, record) => (
          <div>
            <Button style={{marginRight: "8px"}} onClick={() => this.editRelation(record)}>{i18next.t("general:Edit")}</Button>
            <PopconfirmModal
              title={i18next.t("general:Sure to delete") + `: ${record.name} ?`}
              onConfirm={() => this.deleteRelation(record)}
            />
          </div>
        ),
      },
    ];

    return (
      <Card size="small" title={"Channels"}>
        <Table
          rowKey="name"
          loading={this.state.loading}
          dataSource={this.state.relations}
          columns={columns}
          pagination={false}
        />
      </Card>
    );
  }

  render() {
    return (
      <div>
        {this.renderEditor()}
        <div style={{marginTop: "20px"}}>
          {this.renderTable()}
        </div>
      </div>
    );
  }
}

export default ChannelRelationListPage;
