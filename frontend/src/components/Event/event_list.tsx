import { EventListItem } from "./event_list_item";
import * as React from "react";
import { EventModel } from "../../models/event";

export interface EventListProps {
  events: EventModel[];
}

export class EventList extends React.Component<EventListProps, {}> {
  render() {
    const items = this.props.events.map((e) => <EventListItem event={e} />);
    return (
      <table className="table">
        <thead>
          <tr>
            <th>Event</th>
            <th>Where</th>
            <th colSpan={2}>When (start/end)</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>{items}</tbody>
      </table>
    );
  }
}
 