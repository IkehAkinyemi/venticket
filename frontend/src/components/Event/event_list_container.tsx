import * as React from "react";
import { EventList } from "./event_list";
import { EventModel } from "../../models/event";

export interface EventListContainerProps {
  eventListURL: string;
}

export interface EventListContainerState {
  loading: boolean;
  events: EventModel[];
}

export class EventListContainer extends React.Component<
  EventListContainerProps,
  EventListContainerState
> {
  constructor(p: EventListContainerProps) {
    super(p);

    this.state = {
      loading: true,
      events: [],
    };

    fetch(p.eventListURL)
      .then<EventModel[]>((response) => response.json())
      .then((events) => {
        this.setState({
          loading: false,
          events: events,
        });
      });
  }

  render() {
    if (this.state.loading) {
      return <div>Loading...</div>;
    }
    return <EventList events={this.state.events} />;
  }
}
