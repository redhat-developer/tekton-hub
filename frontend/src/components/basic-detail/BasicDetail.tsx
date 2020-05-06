import React, {useState} from 'react';
import {
  Card,
  Flex,
  FlexItem,
  Button,
  Grid,
  Badge,
  GridItem,
  CardHead,
  TextContent,
  Text,
  CardActions,
  ClipboardCopy,
  ClipboardCopyVariant,
  Modal,
  TextVariants,
  DropdownItem,
  DropdownToggle,
  Dropdown,
} from '@patternfly/react-core';
import {
  GithubIcon,
  BuildIcon,
  DomainIcon,
} from '@patternfly/react-icons';
import '@patternfly/react-core/dist/styles/base.css';
import './basicdetail.css';
import {fetchTaskDescription} from '../redux/Actions/TaskActionDescription';
import {connect} from 'react-redux';
import Rating from '../rating/Rating';
export interface BasicDetailPropObject {
  id: any
  name: string;
  description: string;
  rating: number;
  latestVersion: string,
  tags: []
  type: string
  data: []
}

export interface Version {
  version: string,
  description: string,
  rawUrl: string
  webUrl: string
}

export interface BasicDetailProp {
  task: BasicDetailPropObject
  version: Version
}


const BasicDetail: React.FC<BasicDetailProp> = (props: any) => {
  React.useEffect(() => {
    props.fetchTaskDescription(props.version.rawUrl);
    // eslint-disable-next-line
  }, [])

  const taskArr: any = [];
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [descrption, setDescription] = useState(props.task.description);
  const [versions, setVersion] =
    useState(props.task.latestVersion + ' (latest) ');

  const [taskLink, setTaskLink] =
    useState(`kubectl apply -f ${props.version.rawUrl}`);

  const [href, setHref] = useState(`${props.version.webUrl.substring(0,
    props.version.webUrl.lastIndexOf('/') + 1)}`);

  // Dropdown menu to show versions
  const [isOpen, set] = useState(false);
  let dropdownItems: any = [];

  if (props.task.data) {
    fetchTaskDescription(props.version.rawUrl);
    dropdownItems = props.task.data.reverse().map((item: any, index: any) => {
      return <DropdownItem
        key={`res-${item.version}`} id={item.version}
        onClick={version}>{item.version}
      </DropdownItem>;
    });
  }

  // Version for resource
  function version(event: any) {
    props.task.data.forEach((item: any) => {
      if (event.target.text === props.task.latestVersion) {
        setVersion(props.task.latestVersion + ' (latest) ');
      } else {
        setVersion(event.target.text);
      }
      if (event.target.text === item.version) {
        props.fetchTaskDescription(item.rawUrl);

        setHref(`${item.webUrl.substring(0,
          item.webUrl.lastIndexOf('/') + 1)}`);

        setTaskLink(`kubectl apply -f ${item.rawUrl}`);
        setDescription(item.description);
      }
    });
  }

  const ontoggle = (isOpen: React.SetStateAction<boolean>) => set(isOpen);
  const onSelect = () => set(!isOpen);

  // Get tags for resource
  if (props.task.tags != null) {
    props.task.tags.forEach((item: any) => {
      taskArr.push(item.name);
    });
  } else {
    taskArr.push([]);
  }
  // ading icon for details page
  let resourceIcon: React.ReactNode;
  if (props.task.type === 'task') {
    resourceIcon = <BuildIcon
      style={{height: '5em', width: '5em'}} color="#484848" />;
  } else {
    resourceIcon = <DomainIcon
      style={{height: '5em', width: '5em'}} color="#4848484" />;
  }

  return (
    <Flex>

      <Card style={{
        marginLeft: '-2em', marginRight: '-2em',
        marginTop: '-2em', width: '120%', paddingBottom: '2em',
      }}>
        <CardHead style={{paddingTop: '2em'}}>
          <div style={{height: '7em', paddingLeft: '10em', marginTop: '5em'}}>
            {resourceIcon}
          </div>

          <TextContent style={{paddingLeft: '4em', paddingTop: '2em'}}>

            <Text style={{fontSize: '2em'}}>
              {props.task.name.charAt(0).toUpperCase() +
                props.task.name.slice(1)}
            </Text>

            <Text style={{fontSize: '1em'}}>
              <GithubIcon size="md"
                style={{marginRight: '0.5em', marginBottom: '-0.3em'}} />

              <a href={href} target="_">Github</a>
            </Text>

            <Grid>

              <GridItem span={10} style={{paddingBottom: '1.5em'}}>
                {descrption}
              </GridItem>

              <GridItem>
                {
                  taskArr.map((tag: any) => {
                    return (
                      <Badge
                        style={{
                          paddingRight: '1em',
                          marginBottom: '1em', marginRight: '1em',
                        }}
                        key={tag}
                        className="badge">{tag}
                      </Badge>);
                  })
                }
              </GridItem>

            </Grid>

          </TextContent>

          <CardActions style={{marginRight: '3em', paddingTop: '2em'}}>

            <Flex breakpointMods={[{modifier: 'column', breakpoint: 'lg'}]}>
              <FlexItem>
                <Rating />
              </FlexItem>

              <FlexItem style={{marginLeft: '-3em'}}>
                <React.Fragment>
                  {document.queryCommandSupported('copy')}
                  <Button variant="primary"
                    className="button"
                    onClick={() => setIsModalOpen(!isModalOpen)}
                  >
                    Install
                  </Button>

                  <Modal
                    width={'60%'}
                    title={props.task.name.charAt(0).toUpperCase() +
                      props.task.name.slice(1)}
                    isOpen={isModalOpen}
                    onClose={() => setIsModalOpen(!isModalOpen)}
                    isFooterLeftAligned
                  >
                    <hr />
                    <div>

                      <TextContent>
                        <Text component={TextVariants.h2} className="modaltext">
                          Install on Kubernetes
                        </Text>
                        {/* {pipelineLink} */}
                        <Text> Tasks </Text>

                        <ClipboardCopy isReadOnly
                          variant={ClipboardCopyVariant.expansion}>{taskLink}
                        </ClipboardCopy>

                      </TextContent>

                      <br />
                    </div>

                  </Modal>

                </React.Fragment>

              </FlexItem>

              <FlexItem style={{marginLeft: '-2em', marginTop: '1'}}>

                <Dropdown
                  onSelect={onSelect}
                  toggle={<DropdownToggle
                    onToggle={ontoggle}>{versions}
                  </DropdownToggle>}
                  isOpen={isOpen}
                  dropdownItems={dropdownItems}
                />

              </FlexItem>

            </Flex>

          </CardActions>

        </CardHead>

      </Card>

    </Flex >
  );
};

const mapStateToProps = (state: any) => {
  return {
    TaskDescription: state.TaskDescription.TaskDescription,
  };
};

export default connect(mapStateToProps,
  {fetchTaskDescription})(BasicDetail);


