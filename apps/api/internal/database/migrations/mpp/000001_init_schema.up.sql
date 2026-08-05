-- =====================================================
-- MPP SCHEMA - QUEUE MANAGEMENT SYSTEM
-- =====================================================
-- Migration 001: Initialize schema, extensions, and enum types
-- =====================================================

-- Create schema
CREATE SCHEMA IF NOT EXISTS mpp;

-- Enable required extensions (idempotent; also enabled by core)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set search path
SET search_path TO mpp, core, public;

-- =====================================================
-- ENUM TYPES
-- =====================================================

-- Queue ordering mode, configurable per instansi
CREATE TYPE mpp.queue_mode AS ENUM ('FIFO', 'BOOKING_PRIORITY');

-- Counter (loket) operational status
CREATE TYPE mpp.loket_status AS ENUM ('OPEN', 'CLOSED', 'BREAK');

-- How a booking was created
CREATE TYPE mpp.booking_channel AS ENUM ('WHATSAPP', 'WEB');

-- Booking lifecycle (pre-queue)
CREATE TYPE mpp.booking_status AS ENUM ('BOOKED', 'CHECKED_IN', 'EXPIRED', 'CANCELLED');

-- Where a queue item originated
CREATE TYPE mpp.antrian_source AS ENUM ('BOOKING', 'WALK_IN', 'SECOND_SERVICE');

-- Queue item lifecycle (see docs/02-domain/queue-state-machine.md)
CREATE TYPE mpp.antrian_status AS ENUM (
    'WAITING',
    'CALLED',
    'SERVING',
    'HOLD',
    'DONE',
    'SKIPPED',
    'TRANSFERRED',
    'QUEUED_NEXT',
    'CANCELLED'
);

-- Outcome recorded on a serving session
CREATE TYPE mpp.serving_outcome AS ENUM ('DONE', 'SKIPPED', 'TRANSFERRED', 'HOLD');

-- Front-office document verification result
CREATE TYPE mpp.fo_result AS ENUM ('COMPLETE', 'INCOMPLETE');

COMMENT ON SCHEMA mpp IS 'Mall Pelayanan Publik (MPP) queue-management domain';
