/* eslint-disable max-len */
/* eslint-disable react/prop-types */
import React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import {
  Link,
} from 'react-router-dom';
import './index.css';
import {
  Card,
  Badge,
  GalleryItem,
  TextContent,
  CardHead,
  CardHeader,
  CardFooter,
  CardBody,
  CardActions,
  Label,
} from '@patternfly/react-core';
import {OkIcon, DownloadIcon, BuildIcon, DomainIcon} from '@patternfly/react-icons';
export interface TaskPropObject {
  name: string;
  description: string;
  rating: number;
  downloads: number;
  yaml: string;
  tags: [];
}

export interface TaskProp {
  task: TaskPropObject
}

// eslint-disable-next-line
const Task: React.FC<TaskProp> = (props: any) => {
  const tempArr: any = [];
  if (props.task.tags != null) {
    tempArr.push(props.task.tags);
  } else {
    tempArr.push([]);
  }
  let verifiedStatus: any;
  if (props.task.verified === true) {
    verifiedStatus = <div className="vtask" >
      <Label isCompact style={{backgroundColor: '#B8AD8B', fontSize: '0.9em'}}>Verified</Label>
    </div>;
  }
  // for adding icon to task and pipeline
  let resourceIcon: React.ReactNode;
  if (props.task.type === 'task') {
    resourceIcon = <BuildIcon size="xl" color="blue" />;
  } else {
    resourceIcon = <DomainIcon size="xl" color="blue" />;
  };

  return (
    <GalleryItem>
      <Link to={'/detail/' + props.task.id}>
        <Card className="card" isHoverable style={{marginBottom: '2em', borderRadius: '0.5em'}}>
          {verifiedStatus}

          <CardHead>
            <div style={{marginTop: '-2em'}}>
              {resourceIcon}
            </div>

            <CardActions className="cardActions">
              <DownloadIcon style={{marginRight: '0.2em'}} className="download" />
              <TextContent className="text">{props.task.downloads}</TextContent>
              <OkIcon style={{color: 'green'}} />
              <TextContent className="text">{props.task.rating.toFixed(1)}</TextContent>
            </CardActions>
          </CardHead>
          <CardHeader className="catalog-tile-pf-header">
            <span className="task-heading">{props.task.name[0].toUpperCase() + props.task.name.slice(1)}</span>
          </CardHeader>
          <CardBody className="catalog-tile-pf-body">
            <div className="catalog-tile-pf-description">
              <span>
                {`${props.task.description.substring(0, 100)}   ...`}
              </span>
            </div>
          </CardBody>
          <CardFooter className="catalog-tile-pf-footer">
            {
              tempArr[0].map((tag: any) => (<Badge style={{
                marginLeft: '0.2em',
                marginBottom: '1em',
              }} key={tag.Name} className="badge">{tag}</Badge>))
            }
          </CardFooter>
        </Card>
      </Link>
    </GalleryItem>
  );
};
export default Task;
