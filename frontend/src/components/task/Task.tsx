/* eslint-disable max-len */
/* eslint-disable react/prop-types */
import React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
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
} from '@patternfly/react-core';
import {
  StarIcon,
  BuildIcon,
  DomainIcon,
} from '@patternfly/react-icons';
export interface TaskPropObject {
  name: string;
  description: string;
  rating: number;
  downloads: number;
  yaml: string;
  tags: [];
  last_updated_at: string;
  latest_version;
}

export interface TaskProp {
  task: TaskPropObject
}

// eslint-disable-next-line
const Task: React.FC<TaskProp> = (props: any) => {
  const tempArr: any = [];
  if (props.task.tags != null) {
    props.task.tags.forEach((item: any) => {
      tempArr.push(item.name)
    })
  } else {
    tempArr.push([]);
  }

  TimeAgo.addLocale(en)

  // Create relative date/time formatter.
  const timeAgo = new TimeAgo('en-US')

  var catalogDate = new Date(props.task.last_updated_at);

  var diffDays = timeAgo.format(catalogDate.getTime() - 60 * 1000)


  let resourceIcon: React.ReactNode;
  if (props.task.type === 'task') {
    resourceIcon = <BuildIcon size="xl" color="#484848" />;
  } else {
    resourceIcon = <DomainIcon size="xl" color="#484848" />;
  };

  return (
    <GalleryItem>
      <Link to={'/detail/' + props.task.id}>
        <Card className="card" isHoverable style={{marginBottom: '2em', borderRadius: '0.5em'}}>

          <CardHead>
            <div>
              {resourceIcon}
            </div>

            <CardActions className="cardActions">

              <TextContent className="text">v{props.task.latest_version}</TextContent>

              <StarIcon style={{color: '#484848'}} />
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

            <TextContent className="text" style={{marginBottom: "1em", marginLeft: "0.2em"}}>Updated {diffDays} </TextContent>
            {
              tempArr.map((tag: any) => {
                return (
                  <Badge style={{
                    marginLeft: '0.2em',
                    marginBottom: '1em',
                  }} key={`badge-${tag}`}
                    className="badge">{tag}</Badge>
                )
              })
            }
          </CardFooter>
        </Card>
      </Link>
    </GalleryItem>
  );
};
export default Task;
