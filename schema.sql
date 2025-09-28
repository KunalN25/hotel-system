-- HOTELS
CREATE TABLE public.hotels (
  id               integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name             text NOT NULL,
  description      text,
  available_rooms  integer NOT NULL CHECK (available_rooms >= 0),
  total_rooms      integer NOT NULL CHECK (total_rooms >= 0),
  street           text,
  landmark         text,
  locality         text,
  city             text,
  pincode          integer,
  state            text,
  image_urls       text[] DEFAULT ARRAY[]::text[],
  cost_per_night   numeric(10,2) NOT NULL CHECK (cost_per_night >= 0),
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_hotels_city ON public.hotels (city);
CREATE INDEX IF NOT EXISTS idx_hotels_locality ON public.hotels (locality);
CREATE INDEX IF NOT EXISTS idx_hotels_pincode ON public.hotels (pincode);


-- USERS
CREATE TABLE public.users (
  id         integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username   text NOT NULL UNIQUE,
  password   text NOT NULL, -- store password hashes, not plaintext
  created_at timestamptz NOT NULL DEFAULT now()
);


-- BOOKINGS
CREATE TABLE public.bookings (
  booking_id      integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  hotel_id        integer NOT NULL,
  user_id         integer NOT NULL,
  number_of_rooms integer NOT NULL CHECK (number_of_rooms > 0),
  number_of_days  integer NOT NULL CHECK (number_of_days > 0),
  booking_time    timestamptz NOT NULL DEFAULT now(),
  check_in_date   date NOT NULL,
  check_out_date  date NOT NULL,
  status          text NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE public.bookings
  ADD CONSTRAINT fk_bookings_hotels FOREIGN KEY (hotel_id)
    REFERENCES public.hotels (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE public.bookings
  ADD CONSTRAINT fk_bookings_users FOREIGN KEY (user_id)
    REFERENCES public.users (id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON public.bookings (user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_hotel_id ON public.bookings (hotel_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON public.bookings (status);


-- PAYMENTS
CREATE TABLE public.payments (
  id                  integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  booking_id          integer NOT NULL,
  order_id            text,
  amount              numeric(10,2) NOT NULL CHECK (amount >= 0),
  currency            text NOT NULL DEFAULT 'INR',
  status              text NOT NULL, 
  created_at          timestamptz NOT NULL DEFAULT now(),
  checkout_session_id text
);

ALTER TABLE public.payments
  ADD CONSTRAINT fk_payments_bookings FOREIGN KEY (booking_id)
    REFERENCES public.bookings (booking_id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_payments_booking_id ON public.payments (booking_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON public.payments (order_id);


-- IDEMPOTENCY KEYS
CREATE TABLE public.idempotency_keys (
  id               uuid PRIMARY KEY,
  idempotency_key  text NOT NULL,
  user_id          integer NOT NULL,
  endpoint         text NOT NULL,
  request_payload  jsonb,
  response_payload jsonb,
  created_at       timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_idempotency_user_endpoint_key
  ON public.idempotency_keys (user_id, endpoint, idempotency_key);